package config

var templateString = `---
id: {{.ID}}
secret: some-secret
modules:
  crud:
    {{.PrimaryDB}}:
      enabled: true
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
