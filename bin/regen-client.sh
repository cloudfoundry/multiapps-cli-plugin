#!/bin/bash

GREEN_COLOUR='\033[0;32m'
DEFAULT_COLOUR='\033[0m' # No Color

check_dir() {
	if [ -d "$1" ]; then
		echo -e "Folder $1 already exists.\nDo you want to ${GREEN_COLOUR}overwrite${DEFAULT_COLOUR} the files? [y/n]"
		
		read answer
		
		if [ "$answer" != 'y' ] && [ "$answer" != 'Y' ] && [ "$answer" != 'yes' ] && [ "$answer" != 'YES' ]; then
			echo -e "Edit regen folder name and ${GREEN_COLOUR}restart${DEFAULT_COLOUR} script."
			exit 1
		fi
		
		rm -r "$1"
	fi
	
	echo -e "${GREEN_COLOUR}Recreating folder...${DEFAULT_COLOUR}"
	mkdir "$1"
}

go_home=$(echo ~/go)
git_home=$(echo ~/git)
tmp="swagger-regen"
swagger_file="mtarest.yaml"
swagger_file_v2="mtarest_v2.yaml"
client_name="mtaclient"

if [ $# -eq 0 ];
then
    echo "No arguments supplied, generating in dirs ${tmp} and ${tmp}-v2"
elif [ "$1" ]; then
	tmp="$1"
fi

regen_folder="${go_home}/src/github.com/cloudfoundry-incubator/multiapps-cli-plugin/${tmp}"
regen_folder_v2="${regen_folder}-v2"
definition_file="${git_home}/multiapps-controller/com.sap.cloud.lm.sl.cf.api/src/main/resources/${swagger_file}"
definition_file_v2="${git_home}/multiapps-controller/com.sap.cloud.lm.sl.cf.api/src/main/resources/${swagger_file_v2}"

check_dir "${regen_folder}"
check_dir "${regen_folder_v2}"

echo -e "Assuming controller project is under this parent dir: ${GREEN_COLOUR}${git_home}${DEFAULT_COLOUR}"
echo -e "Assuming plugin project is under this parent dir: ${GREEN_COLOUR}${go_home}${DEFAULT_COLOUR}"
echo -e "Reading from\n\t${GREEN_COLOUR}${definition_file}\n\t${definition_file_v2}${DEFAULT_COLOUR}\nGenerating in \n\t${GREEN_COLOUR}${regen_folder}\n\t${regen_folder_v2}${DEFAULT_COLOUR}"

cd "${git_home%/*}"
mvn -f "${git_home}/multiapps-controller/" clean package -DskipTests=true -pl=com.sap.cloud.lm.sl.cf.api

swagger generate client -f ${definition_file} -A http_mta -c ${client_name} -t ${regen_folder}
swagger generate client -f ${definition_file_v2} -A http_mta_v2 -c "${client_name}_v2" -t ${regen_folder_v2}