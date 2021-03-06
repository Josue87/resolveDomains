# resolveDomains

Given a list of domains, you resolve them and get the IP addresses.

# Installation

If you want to make modifications locally and compile it, follow the instructions below:

```
> git clone https://github.com/Josue87/resolveDomains.git
> cd resolveDomains
> go build
```

If you are only interested in using the program:

```
> go get github.com/Josue87/resolveDomains 
```

# Usage

```
> resolveDomains -d domainFiles.txt [-t 150] [-r 8.8.8.8:53]
```

Don't forget the ./ in front of the program name if you are compiling locally!

# Example

![image](https://user-images.githubusercontent.com/16885065/119138781-8bbd9f00-ba42-11eb-87f8-63e38fc93e29.png)

# Author

This code has been developed by:

* **Josué Encinar García** -- [@JosueEncinar](https://twitter.com/JosueEncinar)

Code adapted from:

* [Subdomain_guesser (Black Hat Go)](https://github.com/blackhat-go/bhg/blob/master/ch-5/subdomain_guesser/main.go)

At the request of:

* **Six2dez** -- [@six2dez1](https://twitter.com/six2dez1)

Please feel free to modify the code or use it in your tools.
