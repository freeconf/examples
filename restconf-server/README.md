After you captured your management into a root Node, you'll likely want to serve that mangement capability thru a network aware interface.

Here we add support for RESTCONF to our car application in the `main()` method. There is essentially these steps :
1. Identify the path or paths to your YANG files
2. Create a `Device` object
3. Add the management browser to you application
4. Create a RESTCONF server
5. Apply a startup config

