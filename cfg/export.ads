artifacts builderVersion:"1.1", {

    version "${buildBaseVersion}", {
        group "${groupId}", {
            artifact "${artifactId}", {
                file "${genroot}/out/mta_plugin_linux_amd64", classifier: "linux", extension: "bin"
                file "${genroot}/out/mta_plugin_darwin_amd64", classifier: "darwin", extension: "bin"
                file "${genroot}/out/mta_plugin_windows_amd64.exe", classifier: "windows", extension: "exe"
            }
        }
    }
}
