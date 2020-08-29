package syncman

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spaceuptech/helpers"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// ETCDStore is an object for storing ETCD information
type ETCDStore struct {
	etcdClient                       *clientv3.Client
	kv                               clientv3.KV
	nodeID, clusterID, advertiseAddr string
}

// GetAdminConfig returns the admin config
func (s *ETCDStore) GetAdminConfig(ctx context.Context) (*config.Admin, error) {
	key := fmt.Sprintf("sc/admin-config/%s", s.clusterID)
	res, err := s.kv.Get(ctx, key)
	if err != nil {
		return nil, helpers.Logger.LogError(helpers.GetInternalRequestID(), fmt.Sprintf("Unable to fetch key (%s) from etcd", key), err, map[string]interface{}{})
	}

	if len(res.Kvs) == 0 {
		return nil, helpers.Logger.LogError(helpers.GetInternalRequestID(), fmt.Sprintf("Specified key (%s) doesn't exists in etcd", key), err, map[string]interface{}{})
	}

	// we will get the latest value at index 0, as we haven't specified any revision option in get method
	var cluster *config.Admin
	if err := json.Unmarshal(res.Kvs[0].Value, &cluster); err != nil {
		return nil, helpers.Logger.LogError(helpers.GetInternalRequestID(), "Unable to unmarshal cluster config of etcd to json", err, map[string]interface{}{})
	}

	return cluster, nil
}

type trackedItemMeta struct {
	createRevision int64
	modRevision    int64
	service        *service
	project        *config.Project
}

// NewETCDStore creates new etcd store
func NewETCDStore(nodeID, clusterID, advertiseAddr string) (*ETCDStore, error) {
	config, err := loadConfig()
	if err != nil {
		return &ETCDStore{}, helpers.Logger.LogError(helpers.GetInternalRequestID(), "Unable to load etch config from environment variables", err, nil)
	}

	client, err := clientv3.New(config)
	if err != nil {
		return nil, helpers.Logger.LogError(helpers.GetInternalRequestID(), "Unable to initialize etcd client", err, nil)
	}
	return &ETCDStore{etcdClient: client, nodeID: nodeID, clusterID: clusterID, advertiseAddr: advertiseAddr, kv: clientv3.NewKV(client)}, nil
}

// loadConfig return the config object required for creating etcd store instance, required config is loaded from environment variables
func loadConfig() (clientv3.Config, error) {
	isSSL, err := strconv.ParseBool(os.Getenv("ETCD_HTTP_SSL"))
	if err != nil {
		return clientv3.Config{}, helpers.Logger.LogError(helpers.GetInternalRequestID(), "Unable to parse ETCD_HTTP_SSL environment variable to bool", err, nil)
	}

	endpoints := os.Getenv("ETCD_ENDPOINTS")
	adminUser := os.Getenv("ETCD_USER")
	adminPass := os.Getenv("ETCD_PASSWORD")
	caCert := os.Getenv("ETCD_CACERT")
	publicKey := os.Getenv("ETCD_CLIENT_CERT")
	privateKey := os.Getenv("ETCD_CLIENT_KEY")

	if endpoints != "" {
		return clientv3.Config{}, helpers.Logger.LogError(helpers.GetInternalRequestID(), "ETCD endpoints cannot be empty", err, nil)
	}

	var client clientv3.Config
	if isSSL {
		caCert, err := ioutil.ReadFile(caCert)
		if err != nil {
			return clientv3.Config{}, helpers.Logger.LogError(helpers.GetInternalRequestID(), fmt.Sprintf("Unable to read CA cert from path (%s)", caCert), err, nil)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		// Read the key pair to create certificate
		cert, err := tls.LoadX509KeyPair(publicKey, privateKey)
		if err != nil {
			return clientv3.Config{}, helpers.Logger.LogError(helpers.GetInternalRequestID(), fmt.Sprintf("Unable to read public (%s) or private key (%s) from path", publicKey, privateKey), err, nil)
		}

		client = clientv3.Config{
			TLS: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		}
	}

	client.Endpoints = strings.Split(endpoints, ",")
	client.Username = adminUser
	client.Password = adminPass
	return client, nil
}

// Register registers space cloud to the etcd store
func (s *ETCDStore) Register() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// opts := &api.WriteOptions{}
	// opts = opts.WithContext(ctx)

	lease, err := s.etcdClient.Grant(ctx, 10)
	if err != nil {
		log.Fatal("Could not create a new session with etcd:", err)
	}

	if _, err := s.kv.Put(ctx, fmt.Sprintf("sc/instances/%s/%s", s.clusterID, s.nodeID), s.advertiseAddr, clientv3.WithLease(lease.ID)); err != nil {
		log.Fatal("Could not register space cloud with etcd:", err)
	}

	ch, err := s.etcdClient.KeepAlive(context.Background(), lease.ID)
	if err != nil {
		log.Fatal("Could not renew etcd session:", err)
	}

	go func() {
		for range ch {
			s.Register()
			return
		}
	}()
}

// WatchAdminConfig maintains consistency between all instances of sc
func (s *ETCDStore) WatchAdminConfig(cb func(clusters []*config.Admin)) error {
	// Query all KVs with prefix
	res, err := s.etcdClient.Get(context.Background(), "sc/admin-config/"+s.clusterID, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	clusters := []*config.Admin{
		{
			LicenseKey:   "",
			LicenseValue: "",

			License: "",
		},
	}
	for _, kv := range res.Kvs {
		// Get the id of the item
		if err := json.Unmarshal(kv.Value, clusters[0]); err != nil {
			helpers.Logger.LogInfo(helpers.GetInternalRequestID(), "Sync manager: Could not parse project received", map[string]interface{}{"error": err})
			continue
		}
	}
	cb(clusters)

	ch := s.etcdClient.Watch(context.Background(), fmt.Sprintf("sc/admin-config/%s", s.clusterID), clientv3.WithPrefix())

	go func() {
		for watchResponse := range ch {

			for _, event := range watchResponse.Events {
				if watchResponse.Err() != nil {
					log.Fatal(watchResponse.Err())
				}
				kv := event.Kv

				switch event.Type {
				case mvccpb.PUT:
					if err := json.Unmarshal(kv.Value, clusters[0]); err != nil {
						helpers.Logger.LogInfo(helpers.GetInternalRequestID(), "Sync manager: Could not parse project received", map[string]interface{}{"error": err})
						continue
					}

					cb(clusters)
				}
			}
		}
	}()
	return nil
}

// WatchProjects maintains consistency between all instances of sc
func (s *ETCDStore) WatchProjects(cb func(projects []*config.Project)) error {
	idxID := 3
	itemsMeta := map[string]*trackedItemMeta{}

	// Query all KVs with prefix
	res, err := s.etcdClient.Get(context.Background(), fmt.Sprintf("sc/projects/%s", s.clusterID), clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, kv := range res.Kvs {
		// Get the id of the item
		id := strings.Split(string(kv.Key), "/")[idxID]
		project := new(config.Project)
		if err := json.Unmarshal(kv.Value, project); err != nil {
			helpers.Logger.LogInfo(helpers.GetInternalRequestID(), "Sync manager: Could not parse project received", map[string]interface{}{"error": err})
			continue
		}
		// Store the item
		itemsMeta[id] = &trackedItemMeta{createRevision: kv.CreateRevision, modRevision: kv.ModRevision, project: project}
	}
	cb(s.getProjects(itemsMeta))

	ch := s.etcdClient.Watch(context.Background(), fmt.Sprintf("sc/projects/%s", s.clusterID), clientv3.WithPrefix())

	go func() {
		for watchResponse := range ch {

			for _, event := range watchResponse.Events {
				if watchResponse.Err() != nil {
					log.Fatal(watchResponse.Err())
				}
				kv := event.Kv
				a := strings.Split(string(kv.Key), "/")
				id := a[idxID]
				if a[2] != s.clusterID {
					continue
				}

				switch event.Type {
				case mvccpb.PUT:
					project := new(config.Project)
					if err := json.Unmarshal(kv.Value, project); err != nil {
						helpers.Logger.LogInfo(helpers.GetInternalRequestID(), "Sync manager: Could not parse project received", map[string]interface{}{"error": err})
						continue
					}
					meta, p := itemsMeta[id]
					if !p {
						// AddStateless node if doesn't already exists
						itemsMeta[id] = &trackedItemMeta{createRevision: event.Kv.CreateRevision, modRevision: event.Kv.ModRevision, project: project}
						cb(s.getProjects(itemsMeta))
					}

					// Ignore if incoming create revision is smaller
					if event.Kv.CreateRevision < meta.createRevision {
						break
					}

					// Update if incoming create revision or mod revision is greater
					if event.Kv.CreateRevision > meta.createRevision || event.Kv.ModRevision > meta.modRevision {
						meta.createRevision = event.Kv.CreateRevision
						meta.modRevision = event.Kv.ModRevision
						meta.project = project
						itemsMeta[id] = meta
						cb(s.getProjects(itemsMeta))
					}

				case mvccpb.DELETE:
					meta, p := itemsMeta[id]
					if !p {
						// Ignore if node does not exist
						break
					}

					// Remove if incoming mod revision is greater
					if event.Kv.ModRevision > meta.modRevision {
						// AddStateless node if doesn't already exists
						meta.createRevision = event.Kv.CreateRevision
						meta.modRevision = event.Kv.ModRevision
						delete(itemsMeta, id)
						cb(s.getProjects(itemsMeta))
					}
				}
			}
		}
	}()
	return nil
}

// WatchServices maintains consistency between all instances of sc
func (s *ETCDStore) WatchServices(cb func(scServices)) error {
	idxID := 3
	itemsMeta := map[string]*trackedItemMeta{}

	// Query all KVs with prefix
	res, err := s.etcdClient.Get(context.Background(), fmt.Sprintf("sc/instances/%s", s.clusterID), clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, kv := range res.Kvs {
		// Get the id of the item
		id := strings.Split(string(kv.Key), "/")[idxID]
		service := &service{id: id, addr: string(kv.Value)}
		// Store the item
		itemsMeta[id] = &trackedItemMeta{createRevision: kv.CreateRevision, modRevision: kv.ModRevision, service: service}
	}
	cb(s.getServices(itemsMeta))

	ch := s.etcdClient.Watch(context.Background(), fmt.Sprintf("sc/instances/%s", s.clusterID), clientv3.WithPrefix())

	go func() {
		for watchResponse := range ch {
			for _, event := range watchResponse.Events {
				// Return if watch response contains an error
				if watchResponse.Err() != nil {
					log.Fatal(watchResponse.Err())
				}
				kv := event.Kv
				a := strings.Split(string(kv.Key), "/")
				id := a[idxID]
				if a[2] != s.clusterID {
					continue
				}

				switch event.Type {
				case mvccpb.PUT:
					meta, p := itemsMeta[id]
					if !p {
						// AddStateless node if doesn't already exists
						itemsMeta[id] = &trackedItemMeta{createRevision: event.Kv.CreateRevision, modRevision: event.Kv.ModRevision, service: &service{id: id, addr: string(kv.Value)}}
						cb(s.getServices(itemsMeta))
					}

					// Ignore if incoming create revision is smaller
					if event.Kv.CreateRevision < meta.createRevision {
						break
					}

					// Update if incoming create revision or mod revision is greater
					if event.Kv.CreateRevision > meta.createRevision || event.Kv.ModRevision > meta.modRevision {
						meta.createRevision = event.Kv.CreateRevision
						meta.modRevision = event.Kv.ModRevision
						meta.service = &service{id: id, addr: string(kv.Value)}
						itemsMeta[id] = meta
						cb(s.getServices(itemsMeta))
					}

				case mvccpb.DELETE:
					meta, p := itemsMeta[id]
					if !p {
						// Ignore if node does not exist
						break
					}

					// Remove if incoming mod revision is greater
					if event.Kv.ModRevision > meta.modRevision {
						// AddStateless node if doesn't already exists
						meta.createRevision = event.Kv.CreateRevision
						meta.modRevision = event.Kv.ModRevision
						delete(itemsMeta, id)
						cb(s.getServices(itemsMeta))
					}
				}
			}
		}
	}()

	return nil
}

// SetProject sets the project of the etcd store
func (s *ETCDStore) SetProject(ctx context.Context, project *config.Project) error {
	// todo : why we are setting value as project Id, should'n we marshal the data ?
	_, err := s.kv.Put(ctx, fmt.Sprintf("sc/projects/%s/%s", s.clusterID, project.ID), project.ID)

	return err
}

// SetAdminConfig sets the admin config in the store
func (s *ETCDStore) SetAdminConfig(ctx context.Context, cluster *config.Admin) error {
	data, _ := json.Marshal(cluster)
	_, err := s.kv.Put(ctx, fmt.Sprintf("sc/admin-config/%s", s.clusterID), string(data))

	return err
}

// DeleteProject deletes the project from the etcd store
func (s *ETCDStore) DeleteProject(ctx context.Context, projectID string) error {
	_, err := s.kv.Delete(ctx, fmt.Sprintf("sc/projects/%s/%s", s.clusterID, projectID))
	return err
}

func (s *ETCDStore) getProjects(itemsMeta map[string]*trackedItemMeta) []*config.Project {
	var arrProjects []*config.Project
	for _, item := range itemsMeta {
		arrProjects = append(arrProjects, item.project)
	}
	return arrProjects
}

func (s *ETCDStore) getServices(itemsMeta map[string]*trackedItemMeta) scServices {
	// Sort and store
	var services scServices
	for _, item := range itemsMeta {
		services = append(services, item.service)
	}
	sort.Stable(services)
	return services
}
