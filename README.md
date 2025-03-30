# Proxym

[![Go Reference](https://pkg.go.dev/badge/github.com/nezbut/proxym.svg)](https://pkg.go.dev/github.com/nezbut/proxym)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![CI Status](https://github.com/nezbut/proxym/actions/workflows/ci.yml/badge.svg)](https://github.com/nezbut/proxym/actions)

**Proxym** is a flexible proxy manager for Go applications. It allows you to easily manage the rotation and selection of
proxy servers
through a system of modular strategies. Integrates transparently with the standard `http.Client`.

## Features

✔️ Support for various proxy rotation strategies  
✔️ Easy integration with existing HTTP client  
✔️ Modular architecture for custom strategies  
✔️ Proxy chaining support  

## Install

```bash
go get github.com/nezbut/proxym
```