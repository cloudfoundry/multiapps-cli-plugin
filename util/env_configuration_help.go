package util

const BaseEnvHelpText = `

ENVIRONMENT:
   DEBUG=1                                         Enables the logging of HTTP requests in STDOUT and STDERR.
   MULTIAPPS_CONTROLLER_URL=<URL>                  Overrides the default deploy-service.<system-domain> with a custom URL.
   MULTIAPPS_USER_AGENT_SUFFIX=<STRING>            Appends custom text to User-Agent header. Only alphanumeric, spaces, hyphens, dots, underscores allowed. Max 128 chars, excess truncated.
`
const UploadEnvHelpText = BaseEnvHelpText + `
   MULTIAPPS_UPLOAD_CHUNK_SIZE=<POSITIVE_INTEGER>  Configures chunk size (in MB) for MTAR upload.
   MULTIAPPS_UPLOAD_CHUNKS_SEQUENTIALLY=<BOOLEAN>  Upload chunks sequentially instead of in parallel. By default is false.
   MULTIAPPS_DISABLE_UPLOAD_PROGRESS_BAR=<BOOLEAN> Disable upload progress bar (useful in CI/CD). By default is false.
`
