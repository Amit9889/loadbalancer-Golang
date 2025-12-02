# 🚀 Simple Load Balancer in Go

This project implements a **basic HTTP load balancer** in Go using the `net/http` and `httputil` packages.  
It forwards incoming requests to multiple backend servers using a **Round Robin** strategy.

---

## 📌 Features
- **Reverse Proxy**: Forwards requests to backend servers.
- **Round Robin Algorithm**: Distributes requests evenly across servers.
- **Server Interface**: Abstracts backend
