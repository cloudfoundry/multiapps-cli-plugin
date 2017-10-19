#build scripts must define a class 'build'
import spi
import log

import os
import stat
import sys
import tarfile
import shutil
from distutils.dir_util import copy_tree

class build(spi.BuildPlugin):
  def __init__(self, build_cfg):
    self.build_cfg              =   build_cfg             #it is a good idea to store the build_cfg

    #Environment path variables
    self.go_workspace           =   self.build_cfg.cfg_dir() + "/../gen/go-workspace"
    self.gopath                 =   self.go_workspace + "/go-path"
    self.go_install_dir         =   self.go_workspace + "/go-install-dir"
    self.goroot                 =   self.go_install_dir + "/go"    #set default path
    self.path                   =   self.goroot + "/bin"         #set default path

    #Project parameters
    self.go_project_name        =   "cf-cli-mta-plugin"
    self.go_project_repo        =   "https://github.com/SAP/cf-mta-plugin"
    self.go_project_build_path  =   self.gopath + "/src/" + self.go_project_repo.split(".git",1)[0].split("//")[1]
    self.go_project_version_file=   self.build_cfg.cfg_dir()  + "/VERSION"
    self.go_project_version     =   ""

    #Get imported binary (cfg/import.ais) file path and name
    self.go_binary_file_name    =   "go-binary/go1.8.1.linux-amd64.tar.gz"
    #self.go_binary_file_path	=   self.build_cfg.cfg_dir() + "/../import/content"
    self.go_binary_file_path	=   self.go_project_build_path

    #Build file
    self.build_sh               =   self.build_cfg.cfg_dir() + "/build.sh"

    #Artifacts names
    self.go_project_artifacts   =   ['mta_plugin_linux_amd64', 'mta_plugin_darwin_amd64', 'mta_plugin_windows_amd64.exe']
    self.go_out_path            =    self.build_cfg.cfg_dir() + "/../gen/out"

  def run(self):
    log.info( 'building..' )               #this is where the custom build commands go

    self.clean_workspace()
    self.create_workspace()
    self.get_src_package()
    self.untar_go_binary()
    self.generate_build_sh()
    self.run_build_sh()
    self.copy_binaries_to_gen_out()


  def clean_workspace(self):
    """ This method excute cleaning workspace """
    if ( os.path.isdir(self.go_workspace) ) :
      log.info(' cleaning go workspace path ' + self.go_workspace)
      shutil.rmtree(self.go_workspace)
      log.info(' cleaned go workspace path ' + self.go_workspace)

  def create_workspace(self):
    log.info(' creating environment paths for go ')
    log.info(' creating go workspace path ' + self.go_workspace)
    os.makedirs(self.go_workspace)
    log.info(' creating go path ' + self.gopath)
    os.makedirs(self.gopath)
    log.info(' creating build path ' + self.go_project_build_path)
    os.makedirs(self.go_project_build_path)
    log.info(' creating go install directory path ' + self.go_install_dir)
    os.makedirs(self.go_install_dir)


  def untar_go_binary(self):
    log.info(' untarring go tar file ' + self.go_binary_file_name + '...')
    tar = tarfile.open(self.go_binary_file_path + "/" + self.go_binary_file_name ,"r")
    tar.extractall(self.go_install_dir)
    tar.close()

  def get_src_package(self):
    log.info('Cloning go project repo' + self.go_project_repo)
    os.system("git clone " + self.go_project_repo + " "  + self.go_project_build_path )
    file = open(self.go_project_version_file, "r")
    self.go_project_version = file.read().splitlines()[0]
    file.close()

  def generate_build_sh(self):
    fw = open ( self.build_sh , "w" )
    fw.write("#!/bin/bash -e\n")
    fw.write("export GOPATH=" + self.gopath + "\n")
    fw.write("export GOROOT=" + self.goroot + "\n")
    fw.write("export PATH=$PATH:" + self.path + "\n")
    fw.write("export BUILD_PATH=" + self.go_project_build_path + "\n")
    fw.write("echo \"GOPATH ==> $GOPATH\"\n")
    fw.write("echo \"GOROOT ==> $GOROOT\"\n")
    fw.write("echo \"PATH ==> $PATH\"\n")
    fw.write("echo \"BUILD_PATH ==> $BUILD_PATH\"\n")
    fw.write("cd $BUILD_PATH\n")
    fw.write("./build.sh " + self.go_project_version + "\n")
    fw.close()

    os.system("chmod 777 " + self.build_sh )

  def run_build_sh(self):
    log.info( 'running ' + self.build_sh + '...')
    os.system(self.build_sh)

  def copy_binaries_to_gen_out(self):
    for name in self.go_project_artifacts :
      shutil.copyfile(self.go_project_build_path + "/" + name, self.go_out_path + "/" + name)
    os.system("chmod 777 " + self.go_out_path + "/*")
    os.chdir(self.go_out_path)
    os.chdir(self.build_cfg.cfg_dir() + "/..")

  #set groupID and artifacts names which is used by export.ads
  def deploy_variables(self):
    return {'groupId'       : "com.SAP.golang", 'artifactId'    : "cf-cli-mta-plugin"}
####
