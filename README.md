### LittleEcom
Just another self-project. It's my hobby

### How it works
It's microservices based. Here's a picture explaining overall architecture
![app_architecture.png](assets/app_architecture.png)

### Request flow
![request_flow.png](assets/request_flow.png)

### About dumper.sh
It will continually get logs from every deployment from every pod, dumping each
deployment log in a different file. This is just a workaround, so I can debug
without having to rely on complex logging tools.