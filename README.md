# Bark
Bark is supposed to be a small library and easy-to-use library that uses PostgreSQL for collecting logs from multiple sources. It is also supposed to contain a web server which can accept the logs using REST calls.

## Introduction to the problem
When we start off with smaller projects, logging is not an issue. However, as we grow in size and start creating multiple services and start logging more things, logging becomes more and more problematic. Now, there are pretty great projects, both open and properietary out there that make log collection, search and analytics on terabyte scale possible. However, between the basic HTTP server and the enterprise scale, as we grow, we still need to be able to store, process, analyze and search through our logs. While plaintext logs do make sense, they are not great when you are trying to filter out events, especially in a multi-service installation.

# Aim of this project
Bark aims to fill the gap that exists between simple file-based logging and a large log aggregation and analytics solution by using an already-well-known technology (PostgreSQL) as a log storage server. It is worth noting that PostgreSQL is not great at being a great Log server by design but is just good enough for mid-size installations (say about a dozen services) that has not gone global yet. So bark is not aimed at becoming the most powerful and flexible log aggregation service but useful enough for searching and filtering logs.

We want this service to be available as a web service which can be called by other services to send in the log entries. Bark also aims to be a library which can be included in your project and be called from within the project as a function.

# Language of Choice
It has to be written in golang for the sake of being great at handling incoming traffic bursts, lining up the logs and sending them in a load-controlled manner to PostgreSQL (because PostgreSQL is not optimized for that kinda stuff).

# What is it NOT?
- It is not a replacement for Plaintext logs - Bark should be able to write to a plaintext log file in parallel to throwing items into Postgres.
- It is not a replacement for an APM - We don't want to throw in Application uptime or Performance Monitoring. Bark is not supposed to a monitoring solution at all.
- It is not trying to replace ELK or NewRelic or Datadog or Fluentd or any log analytics system that is dedicated towards the 100 GB+ scale of logs.
- It is not a CLI tool or a Web server at this point - we don't want to start off with a CLI tool or a Web Service to view your logs or filter them. You want that, go run a query against PostgreSQL directly using your terminal or GUI tool of your choice!

Social Media Preview image Photo by [Lora Ninova](https://unsplash.com/@lorannva?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText) on [Unsplash](https://unsplash.com/photos/U86FnrpRR0k?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText)
  
