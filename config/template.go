package config

/**
 * Contains templates for config files
 */

const templateString = `---
id: ID
secret: some-secret
modules:
  crud:
    primary:
      conn: CONN
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
    broker: localhost
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
`//-- end const templateString

const defaultTemplate = templateString

