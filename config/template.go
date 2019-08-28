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
    name: {{.Name}}
    secret: some-secret
    modules:
      crud:
        {{.PrimaryDB}}:
          enabled: true
          conn: {{.Conn}}
          collections:
            default:
              isRealtimeEnabled: true
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
        enabled: true
        broker: nats
        conn: nats://localhost:4222
        services:
          default:
            functions: 
              default:
                rule: 
                  rule: allow
      realtime:
        enabled: true
        broker: nats
        conn: nats://localhost:4222
      pubsub:
        enabled: false
        broker: nats
        conn: nats://localhost:4222
        rules:
          - subject: /
            rule:
              publish:
                rule: allow
              subscribe:
                rule: allow
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
gateway: 
  routes: []
  internalRoutes: []
`
