# latest-jdk

Downloads the URL of the latest JDK, using the Adoptium API.

You can use commands like the following to *download* the latest JDK,
assuming `latest-jdk` is in your PATH:

```bash
curl -LO $(latest-jdk)
curl -LO $(latest-jdk -lts)
curl -LO $(latest-jdk -release 11)
curl -LO $(latest-jdk -os linux -arch x64 -release 18)
```

The program will use your current OS and architecture
as its defaults.

You can type `latest-jdk -h` for help/usage information, including
information about other options.

## Build

Use `go build` to build.
