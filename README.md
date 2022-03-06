# latest-jdk

Downloads the URL of the latest JDK, using the Adoptium API.

Use `go build` to build.

You can use commands like the following to *download* the latest JDK,
assuming `latest-jdk` is in your PATH:

```bash
curl -LO $(latest-jdk)
curl -LO $(latest-jdk -lts)
curl -LO $(latest-jdk -os linux -arch x64 -release 16)
```

You can type `latest-jdk -h` for help/usage information, including
information about other options.
