# Storygenie Backend

This is the backend for the Storygenie app. It is a Go application that uses the Gin framework. It is a REST API that is used by the Storygenie app to store and retrieve data.

## Database
The database provider of choice is [Planetscale](https://planetscale.com/), a MySQL database plattform based on [Vitess](https://vitess.io/).
This decision was made due to their generous free tier and ability to scale through [horizontal sharding](https://planetscale.com/sharding).
You can replace this with any MySQL database and set it in the `.env` file though.

## Setup

To run the application, you need to have Go installed. You can download it from [here](https://golang.org/dl/).
Then, you need to setup local environment variables. You can do this by copying the `.env.example` file to `.env` and filling in the values.

## Generate types from OpenAPI spec

To generate the types from the OpenAPI spec, we use [oapi-codegen](https://github.com/deepmap/oapi-codegen). To generate the types, run the following command:

```bash
go generate
```

## Operations

This project uses gcloud as cloud service. To deploy, obtain a valid serviceAccount for a project and replace the `foobar` occurences below with your project name, then run:

Build:

```bash
gcloud builds submit --tag us-central1-docker.pkg.dev/foobar/cloud-run-source-deploy/storygeniebackend:latest
```

Deploy:

```bash
gcloud run deploy storygeniebackend \
--image=us-central1-docker.pkg.dev/foobar/cloud-run-source-deploy/storygeniebackend:latest \
--region=us-central1 \
--project=foobar \
 && gcloud run services update-traffic storygeniebackend --to-latest --region=us-central1 --project=foobar
```
