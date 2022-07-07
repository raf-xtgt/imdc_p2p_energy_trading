# imdc_p2p_energy_trading



## About the project
With Malaysia's growing population and improved lifestyle, energy demand will only increase posing new challenges and problems to the energy sector. We propose a peer-to-peer (P2P) web-based solar energy trading web application with the aim of optimizing energy usage by homes, such that those with higher electricity demand, can trade energy instead of paying higher bills. Thus also providing a second income stream for prosumers for their production of the tradable renewable solar energy. We use Machine Learning to predict energy consumption and production. Based on these predictions, users can buy and sell renewable energy. To ensure maximum social welfare, we implement a Double Auction Mechanism (DAM) to optimally allocate tradable energy between prosumers where they are guaranteed a fiat reward. The transactions are stored on a Proof-of-Authority (PoA) Distributed Ledger. The back end is built using Go programming language due to its speed, easy syntax and its increasing popularity as a server side language. The Angular framework is used for the front end development with Typescript used for middleware service implementation. Energy forecasting and DAM are written using Python libraries. Moreover, the platform uses the MongoDB NoSQL cloud database to securely and reliably store data. The Locust stress testing tool is used for performance benchmarking.


## Pre-requisites
#### Install Angular13
```
npm install -g @angular/cli
```

### Angular Material for Frontend
```
ng add @angular/material
```

### Install Bootstrap to assist in UI development
```
npm install bootstrap
```

###  [Install Node.js following the instructions](https://phoenixnsp.com/kb/install-node-js-npm-on-windows)

### Install Node.js Package Manager
```
npm install
```

## Run the web application

### Backend

```
cd backend/src
```

```
go run *.go
```

### Frontend

```
cd frontend/web-app/
```
```
ng serve
```

Open browser port

```
localhost:4200
```

### Project Demonstration Video
[Demo](https://www.youtube.com/watch?v=8qXSCftU9V8)

