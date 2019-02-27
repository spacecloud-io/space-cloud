package config

var templateString = `---
id: {{.ID}}
secret: some-secret
modules:
  crud:
    {{.PrimaryDB}}:
      conn: {{.Conn}}
      isPrimary: true
      collections:
        users:
          isRealtimeEnabled: false
          rules:
            create:
              rule: allow
            read:
              rule: allow
            update:
              rule: allow
            delete:
              rule: allow
  auth:
    email:
      enabled: false
  faas:
    enabled: false
    nats: nats://localhost:4222
  realtime:
    enabled: false
    kafka: localhost
  fileStore:
    enabled: false
    storeType: local
    conn: ./
    rules:
      rule1:
        prefix: /
        rule:
          create:
            rule: allow
          read:
            rule: allow
          delete:
            rule: allow
`
