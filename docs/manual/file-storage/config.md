# Configuring the file storage module

The config pertaining to file management module can be found inside the `fileStore` key under the `modules` object. Here's the snippet:

```yaml
modules:
  fileStore:
    enabled: true
    storeType: local
    conn: /tmp/data
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

  # Config for other modules go here
```

Inside `fileStore` you have to speicify the following fields:
- **enable:** Setting this field to `true` enables the file storage module.
- **storeType:** Choose the storage technology to back the file storage module. The possible options are `local` and `amazon-s3`.
- **conn:** Connection is the region name for `amazon-s3` and the root file path for `local`.
- **rules:** Are the [security rules](/docs/security/file-storage) used to secure access to the file storage module. It's essentially a map which contains the `rule` and the `prefix` on which to apply the rule. 

The snippet shown above configues Space Cloud to use the local filesystem (`local`). The files will be stored under the `/tmp/data` directory. There is a single rule named `rule1` which allows all file storage operations (`create`, `read` and `delete`). All rules are applied based on prefix matching on the URL.

Here's how you can add different rules for different prefixes.

```yaml
modules:
  fileStore:
    enabled: true
    storeType: local
    conn: /tmp/data
    rules:
      imageRule:
        prefix: /images/:userId
        rule:
          create:
            rule: allow
          read:
            rule: allow
          delete:
            rule: allow
      musicRule:
        prefix: /music
        rule:
          create:
            rule: allow
          read:
            rule: allow
          delete:
            rule: allow

  # Config for other modules go here
```

> **Note:** The `allow` rule must only be used cautiously. It must never be used for create and delete operations in production. You can read more about security rules for the files module [here](/docs/security/file-storage). 

## Next steps

Now you know the basics of the file storage module. The next step would be checking out the functions module to extend Space Cloud by writing custom logic on backend!

<div class="btns-wrapper">
  <a href="/docs/file-storage/delete" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/functions/overview" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>