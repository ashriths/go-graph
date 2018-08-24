# Go-Graph - 
## A fault-tolerant, scalable distributed graph database written in golang

Today, there are a breadth of applications that process highly connected data by exploiting the connectivity and relationship between bits of information. Social Networking, recommendations based on relationships, Knowledge graphs to improve information retrieval, fraud detection etc are a few use cases for highly connected data. Using relational databases to build such applications poses a few challanges. SQL is not designed to pro- cess relationships well. As the number of relationships grows, the SQL query becomes quite complicated and suffers from performance issues. With highly connected data, the types of information that need to be processed change rapidly. This requires the database system to be flexible to accommodate new data types. Relational Databases don’t evolve their schema very well; they are relatively inflexible. On the other hand, Graph databases are optimized for efficient storage and retrieval of highly connected data and are hence, key to building applications that look to leverage connectivity among bits of information. The need for a distributed, fault-tolerant, scalable graph database is imminent. However, designing a distributed and sharded graph database has proved to be a challenging task due to lack of record-wise separation and strong connections in graph data. Go-graph is a distributed, fault tolerant graph database which handles the intricacies of storing and querying graph data.


## Features
1. Consistent - Graph traversals from multiple vertices which may reside in different machines should still result in the seeing the same graph structure and properties at any point of time after an allowed amount of delay or till the requested version number propagates. Thus, its eventually consistent.
2. Fault tolerant and Replicated - System will sustain node failures at all times
3. Sharded - Not all graph data is replicated everywhere. This can also be extended to geo- local sharding to improve data locality based on query analytics

![Architecture Design](/paper/comp.png "Component Interaction Design")


[Read more about this project and learn about benchmarks and evaluation](README.pdf)
