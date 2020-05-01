## Architecture

The TF-Operator consists of a Stack CRD which maintains the specification and the status of the infrastructure deployed with Terraform. The Spec consists of a reference to a ConfigMap with the tf files that define the infrastructure, and a reference to a Secrect with the tfvars for an specific deployment of this infrastructure (for example, the number of server instances to be deloyed). The ConfigMap is immutable, while the Secret with the tfvars can be modified. The Status of the stack is formed by a Secret with the tfstate and a field with the tfout embedded as a string.

The TF-Operator watches the tfvars and triggers and triggers a Job to run a `terraform apply` command. The Job mounts the Configmap as a directory and the tfvars and state scretes. On finalization, the Job updates the tfstate Secret and the tfout in the Stack status section.

When the Stack CRD is deleted, the TF-Operator launches a Job to exectue a `terraform destroy` command, mounting the configuration, tfvars and current state.

```
                                               +------------------+
                                               | Multiple tf      |
                                          +----+ files (immutable)|
   Stack                                 /     +------------------+
+-----------+                       +---+-------+
|   Spec    |                 +---->|ConfigMap +---+
+-----------+                 |     +----------+   |    +-------------+
|  Config   +-----------------+                    |    | tfvars file |
+-----------+                              +-------·----+ (mutable)   |
|  Tfvars   +------------------+          /        |    +-------------+
+-----------+                  |     +---+-----+   |    
|  Status   |                  +---->+ Secret  +-->+
+-----------+     +---------+        +----+----+   |
|  tfstate  +---->+ Secret  +-----+       ^        |
+-----------+     +----+----+     |       |        |
|  tfout    |          ^          +-------·------->+
| (embedded)|          |                  |        |
+-----------+          |  +--------+      |        |
      ^                |  |Operator|      |        |
      |                |  |        | Watch|        |
      +----------------+  |        +------+        |
                       |  |        |               |
                       |  +---+----+               |
                       |      |                    |
                       |      | Launch             |
                       |      |                    |
                       |      v                    |
                       |  +---+----+               |
                       +--+  Job   |     Mount     |
                    Update|        +<--------------+
                          |        |
                          |        |
                          +--------+
```
