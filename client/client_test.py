import etcd
import sys
import yaml

USER = "user"
PROJECT = "project"
CONFIG = "sample-hivy.yaml"
CONFIG_LOC = "./"

CLIENT = etcd.Etcd()

def SetConfigRow(key,value):
    tab = value.split()
    for i,elem in enumerate(tab):
        #set USER/PROJECT/config/<key>/<index> -d value=<value>
        extkey = USER+"/"+PROJECT+"/"+"config"+"/"+key+"/"+str(i)
	print CLIENT.set(extkey,elem)

def Set(key,value):
    extkey = USER+"/"+PROJECT+"/"+key
    print CLIENT.set(extkey,value)

def Up():
    Set("state","run")

def upper(counter,tab):
    if not(counter<len(tab)-1):
        tab.append("")
        return len(tab)-1 , tab
    return counter+1 , tab

#TODO save user/project input as default
#lol, ugly as ....
if __name__ == '__main__':
    args = sys.argv[1:]
    
    select = -1
    select,args = upper(select,args)
    if args[select] == "setconfig":
	select,args=upper(select,args)
	if args[select] != "":
            USER = args[select]
            select,args=upper(select,args)
            if args[select] != "":
                PROJECT = args[select]
        stram = open(CONFIG_LOC+CONFIG,"r")
        docs = yaml.load(stram)
        #cmd
        for k,v in docs.items():
            SetConfigRow(k,v)
            Set("state","loaded")

    select = -1
    select,args = upper(select,args)
    if args[select] == "up" :
	select,args=upper(select,args)
	if args[select] != "":
            USER = args[select]
            select,args=upper(select,args)
            if args[select] != "":
                PROJECT = args[select]
        #cmd
        Up()
