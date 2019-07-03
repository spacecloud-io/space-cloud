package config

var templateString = `---
ssl:
  enabled: false
admin:
  user: {{.AdminName}}
  pass: {{.AdminPass}}
  role: {{.AdminRole}}
  secret: {{.AdminSecret}}
projects:
  - id: {{.ID}}
    secret: some-secret
    modules:
      crud:
        {{.PrimaryDB}}:
          enabled: true
          conn: {{.Conn}}
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
      functions:
        enabled: false
        broker: nats
        conn: nats://localhost:4222
        rules:
          service1:
            function1:
              rule: allow
      realtime:
        enabled: false
        broker: nats
        conn: nats://localhost:4222
      fileStore:
        enabled: false
        storeType: local
        conn: ./
        rules:
          - prefix: /
            rule:
              create:
                rule: allow
              read:
                rule: allow
              delete:
                rule: allow
      static:
        enabled: false
        routes:
        - prefix: /
          path: ./public
`
var templateStringMissionControl = `---
id: {{.ID}}
secret: some-secret
ssl:
  enabled: false
modules:
  crud:
    {{.PrimaryDB}}:
      conn: {{.Conn}}
      enabled: true
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
  functions:
    enabled: false
    broker: nats
    conn: nats://localhost:4222
    rules:
      service1:
        function1:
          rule: allow
  realtime:
    enabled: false
    broker: nats
    conn: nats://localhost:4222
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
  static:
    enabled: true
    routes:
    - prefix: /mission-control
      path: {{.HomeDir}}/.space-cloud/mission-control-v{{.BuildVersion}}/build
`
