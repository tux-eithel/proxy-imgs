# proxy-imgs
Proxy specific path (like images folder) to another server. Made in Go (golang)

###Use case
You have just synced your repository, and your DB (from production environment) of a WordPress site. Your local site is running in a vm.  
If you navigate the site, you will see the site "broken" because all the images inside `wp-content/uploads` are not in the repository and they may are expensive to sync (GB and GB of data)

Using this little script you can proxy specific path to another site

###Install
`go get -u github.com/tux-eithel/proxy-imgs`

###Example
Your local website is iside a VM and have this data

```
ADDRESS: http://local-site.dev
IP: 192.168.2.10

```
Your remote site has this data
```
ADDRESS: http://remote-site.org
IP: xxx.xxx.xxx.xxx
```
Your Machine IP is `192.168.1.2` and your `/etc/hosts` file has this entry `192.168.1.2    local-site.dev`

You can proxy spcific path using this command:
```
./proxy-imgs -p 8080 -o "http://192.168.2.10:80" -r "http://remote-site.org:80" -f ""wp-content/uploads/*"
```
and navigate your local site browsing http://local-site.dev:8080.


Now all the traffic that NOT match the  `-f` regular expression will be forwarded to the VM (when local-site.dev is running) meanwhile url which match the `-f` regular expression will be forwarded to orgin.-ite.dev

###Tips
You have to use the `ip address` with `-o` flag because Wordpress hard coded domain in urls.

You can specify multiple pattern `-f "*.jpg" -f "*.pdf"`

If you use the `-s` flag, script checks if remote site (`http://remote-site.org` in this example) respond with a status code > 400. If true, request will be forwared to `http://local-site.dev`

###Help
```
Usage of ./proxy-imgs:
  -f value
    	Patter to proxy to Remote site: -f ".jpg?" -f "wp-content/uploads/*"
  -o string
    	Origin site: -o "http://site1.dev:80"
  -p int
    	Port where bind the service -p 8081 bind service to port (default 80)
  -r string
    	Remote site: -r "http://site2.dev:80"
  -s	if remote site gives error code > 400, request will be proxy to Origin. Default FALSE

```