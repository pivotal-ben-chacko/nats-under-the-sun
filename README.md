   

# NATs Inspector

NATs Inspector is a tool that can be used to view messages being sent from the director VM and any VM that has been deployed by the director.

## Connecting to Director NATs

In order to connect to Bosh NATs server we need to establish trust using certificates and private keys, and we also need to understand some of the permissions around what you are able to subscribe to and publish to.

*Director NATs Config*

    DIRECTOR_PERMISSIONS: {
        publish: [
          "agent.*",
          "hm.director.alert"
        ]
        subscribe: ["director.>"]
      }

Above the director is given permission to subscribe to any subject that begins with “director.>” the ‘>’ is a wildcard indicator. In addition the director is allowed to publish to any subject that starts with “agent.*”.  

  

    AGENT_PERMISSIONS: {
        publish: [
          "hm.agent.heartbeat._CLIENT_ID",
          "hm.agent.alert._CLIENT_ID",
          "hm.agent.shutdown._CLIENT_ID",
          "director.*._CLIENT_ID.*"
        ]
        subscribe: ["agent._CLIENT_ID"]
      }

Above are the permissions given to the agent.


### Required certificates:

We use the [nats-io](https://github.com/nats-io/nats.go) go library to connect to the NATs server on the director. This library provides the following method to accomplish this:

    nats.Connect("nats://SERVER-IP-ADDRESS:4222",
    nats.ClientCert("./cert.pem", "./private.key"),
    nats.RootCAs("./root_ca_cert.pem"))

