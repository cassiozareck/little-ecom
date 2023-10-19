## Old things: not used anymore
Old efk stack was handled through k8s manifest, now its just helm; 

## Logs
older logs were emitted through a dumper who could get these logs without efk

## Tcp middleware
I was trying to set middleware in each pod who could catch tcp packets, with the purpose of logs and in-transit packet operation. 
But it was very heavy and sometimes the redirection to the real expected container did not work properly 
