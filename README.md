# adoptium-jdk

Downloads the URL of the latest Adoptium JDK, using the Adoptium API.

You can use commands like the following to *download* the latest Adoptium JDK,
assuming `adoptium-jdk` is in your PATH:

```bash
curl -LO $(adoptium-jdk)
curl -LO $(adoptium-jdk -lts)
curl -LO $(adoptium-jdk -lts -jv)  # print the version only
curl -LO $(adoptium-jdk -release 11)
curl -LO $(adoptium-jdk -os linux -arch x64 -release 18)
```

The program will use your current OS and architecture as its defaults.

You can type `adoptium-jdk -h` for help/usage information, including information
about other options.

The `scripts` directory has scripts that use `adoptium-jdk` and upgrade the JDK.
You can use these scripts as a starting point to upgrade your own setup.

## Build

Use `go build` to build.
