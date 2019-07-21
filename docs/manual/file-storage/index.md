# File storage

The file storage module provides REST APIs for file management. You can select the storage technology for the file management module. Each project can work with only one storage technology. Currently supported storage technologies are:
- Amazon S3
- Local Filesystem

You can currently do the following operations in the file management module:
- Create a folder
- Upload a file
- List all files / folders within a folder
- Delete a file / folder

## Enable the file storage module

The config pertaining to file management module can be found inside the `fileStore` key under the `modules` object. Here's the snippet:

```yaml
modules:
  fileStore:
    enabled: true
    storeType: local
    conn: /tmp/data
    rules:
      - prefix: /
        rule:
          create:
            rule: allow
          read:
            rule: allow
          delete:
            rule: allow

  # Config for other modules go here
```

The snippet shown above configures Space Cloud to use the local filesystem (`local`). The files will be stored under the `/tmp/data` directory. There is a single rule named `rule1` which allows all file storage operations (`create`, `read` and `delete`). All rules are applied based on prefix matching on the URL.

You can learn more about on the various parameters available for configuring the file storage module [here](/docs/file-storage/config).


## Next steps

Now you know the basics of the file storage module. The next step would be diving deeper into the [configuration](/docs/file-store/config) and its structure. You could also continue to see how to use the file storage module on the frontend.

<div class="btns-wrapper">
  <a href="/docs/realtime/overview" class="waves-effect waves-light btn primary-btn-border btn-small">
    <i class="material-icons btn-with-icon">arrow_back</i>Previous
  </a>
  <a href="/docs/file-storage/upload-file" class="waves-effect waves-light btn primary-btn-fill btn-small">
    Next<i class="material-icons btn-with-icon">arrow_forward</i>
  </a>
</div>