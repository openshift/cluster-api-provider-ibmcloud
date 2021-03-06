??    f      L  ?   |      ?  z   ?  ?   	  <   ?	  S   
  <   b
  c  ?
  ?    .   ?  "   ?  4   
     ?     \    {  X   ?  o   ?    J  v   L  t   ?  ?  8  ;   ?  [   9  J   ?  a   ?  ?   B  ?      ?   ?  %   u  W   ?     ?  u     4   ?  -   ?  3   ?  2        Q  *   e  .   ?  *   ?  0   ?  0     0   L  "   }     ?  *   ?  A   ?     +  )   I     s     ?      ?  (   ?     ?  `     ?   m  ?   	     ?     ?  $   ?     ?       a   0  s   ?  B     +   I  +   u  6   ?  q   ?  /   J   1   z   '   ?      ?   &   ?   %   !  (   :!  #   c!      ?!     ?!  9   ?!     "      "  #   :"  ?   ^"  H   ?"  &   *#  e   Q#  ?   ?#  E   ?$  a   ?$  ?   E%  ?   &     ?&     ?&  =   '  $   T'     y'  &   ?'  +   ?'     ?'  r   (     t(  /   ?(  ?  ?(  ?   Y*  ?   ?*  A   ?+  P   ?+  =   6,  ?  t,  ?   .  /   ?/  (   *0  5   S0  +   ?0  1   ?0  '  ?0  ^   2  u   n2  )  ?2  ?   4  }   ?4  ?  5  ;   	7  p   E7  ]   ?7  l   8  ?   ?8  ?   <9  ?   :  .   ?:  r   ?:  $   _;  ?   ?;  :   <  3   F<  =   z<  0   ?<     ?<  )   ?<  3   &=  4   Z=  5   ?=  -   ?=  -   ?=  +   !>     M>  (   i>  P   ?>     ?>  ,   ?  !   .?     P?  #   m?  (   ??     ??  e   ??  ?   <@  ?   ?@  )   ?A  -   ?A  $   ?A  &   B     AB  e   \B  ?   ?B  M   VC  3   ?C  7   ?C  D   D  ?   UD  9   ?D  1   E  $   FE     kE  $   ?E  +   ?E  *   ?E  +   F  0   3F  %   dF  A   ?F     ?F     ?F  (   G  ?   +G  P   ?G  #   H  i   0H    ?H  P   ?I  X   ?I  ?   HJ  ?   6K     /L     IL  C   gL     ?L     ?L  $   ?L  4   M      CM  ?   dM     ?M  2   ?M     	   H       -   #                  3      `       d                  C          I       A          1           >   0          !           "   (       L   %       5   J   ?   4   )   b   Z   @      f   F       =         ;              c   ^         9   [   e   M      a   ,      S   '      \          Q             .   V   T   W       B      Y          E      6      :      X   &   /       P       D   K      U   2                      _   7   ]   <   8           R          $      G              O   N   *   
   +    
		  # Show metrics for all nodes
		  kubectl top node

		  # Show metrics for a given node
		  kubectl top node NODE_NAME 
		# Get the documentation of the resource and its fields
		kubectl explain pods

		# Get the documentation of a specific field of a resource
		kubectl explain pods.spec.containers 
		# Print flags inherited by all commands
		kubectl options 
		# Print the client and server versions for the current context
		kubectl version 
		# Print the supported API versions
		kubectl api-versions 
		# Show metrics for all pods in the default namespace
		kubectl top pod

		# Show metrics for all pods in the given namespace
		kubectl top pod --namespace=NAMESPACE

		# Show metrics for a given pod and its containers
		kubectl top pod POD_NAME --containers

		# Show metrics for the pods defined by label name=myLabel
		kubectl top pod -l name=myLabel 
		Convert config files between different API versions. Both YAML
		and JSON formats are accepted.

		The command takes filename, directory, or URL as input, and convert it into format
		of version specified by --output-version flag. If target version is not specified or
		not supported, convert to latest version.

		The default output will be printed to stdout in YAML format. One can use -o option
		to change to output destination. 
		Create a namespace with the specified name. 
		Create a role with single rule. 
		Create a service account with the specified name. 
		Mark node as schedulable. 
		Mark node as unschedulable. 
		Set the latest last-applied-configuration annotations by setting it to match the contents of a file.
		This results in the last-applied-configuration being updated as though 'kubectl apply -f <file>' was run,
		without updating any other parts of the object. 
	  # Create a new namespace named my-namespace
	  kubectl create namespace my-namespace 
	  # Create a new service account named my-service-account
	  kubectl create serviceaccount my-service-account 
	Create an ExternalName service with the specified name.

	ExternalName service references to an external DNS address instead of
	only pods, which will allow application authors to reference services
	that exist off platform, on other clusters, or locally. 
	Help provides help for any command in the application.
	Simply type kubectl help [path to command] for full details. 
    # Create a new LoadBalancer service named my-lbs
    kubectl create service loadbalancer my-lbs --tcp=5678:8080 
    # Dump current cluster state to stdout
    kubectl cluster-info dump

    # Dump current cluster state to /path/to/cluster-state
    kubectl cluster-info dump --output-directory=/path/to/cluster-state

    # Dump all namespaces to stdout
    kubectl cluster-info dump --all-namespaces

    # Dump a set of namespaces to /path/to/cluster-state
    kubectl cluster-info dump --namespaces default,kube-system --output-directory=/path/to/cluster-state 
    Create a LoadBalancer service with the specified name. A comma-delimited set of quota scopes that must all match each object tracked by the quota. A comma-delimited set of resource=quantity pairs that define a hard limit. A label selector to use for this budget. Only equality-based selector requirements are supported. A label selector to use for this service. Only equality-based selector requirements are supported. If empty (the default) infer the selector from the replication controller or replica set.) Additional external IP address (not managed by Kubernetes) to accept for the service. If this IP is routed to a node, the service can be accessed by this IP in addition to its generated service IP. An inline JSON override for the generated object. If this is non-empty, it is used to override the generated object. Requires that the object supply a valid apiVersion field. Approve a certificate signing request Assign your own ClusterIP or set to 'None' for a 'headless' service (no loadbalancing). Attach to a running container ClusterIP to be assigned to the service. Leave empty to auto-allocate, or set to 'None' to create a headless service. ClusterRole this ClusterRoleBinding should reference ClusterRole this RoleBinding should reference Convert config files between different API versions Copy files and directories to and from containers. Create a TLS secret Create a namespace with the specified name Create a secret for use with a Docker registry Create a secret using specified subcommand Create a service account with the specified name Delete the specified cluster from the kubeconfig Delete the specified context from the kubeconfig Deny a certificate signing request Describe one or many contexts Display clusters defined in the kubeconfig Display merged kubeconfig settings or a specified kubeconfig file Display one or many resources Drain node in preparation for maintenance Edit a resource on the server Email for Docker registry Execute a command in a container Forward one or more local ports to a pod Help about any command If non-empty, set the session affinity for the service to this; legal values: 'None', 'ClientIP' If non-empty, the annotation update will only succeed if this is the current resource-version for the object. Only valid when specifying a single resource. If non-empty, the labels update will only succeed if this is the current resource-version for the object. Only valid when specifying a single resource. Mark node as schedulable Mark node as unschedulable Mark the provided resource as paused Modify certificate resources. Modify kubeconfig files Name or number for the port on the container that the service should direct traffic to. Optional. Only return logs after a specific date (RFC3339). Defaults to all logs. Only one of since-time / since may be used. Output shell completion code for the specified shell (bash or zsh) Password for Docker registry authentication Path to PEM encoded public key certificate. Path to private key associated with given certificate. Precondition for resource version. Requires that the current resource version match this value in order to scale. Print the client and server version information Print the list of flags inherited by all commands Print the logs for a container in a pod Resume a paused resource Role this RoleBinding should reference Run a particular image on the cluster Run a proxy to the Kubernetes API server Server location for Docker registry Set specific features on objects Set the selector on a resource Show details of a specific resource or group of resources Show the status of the rollout Synonym for --target-port The image for the container to run. The image pull policy for the container. If left empty, this value will not be specified by the client and defaulted by the server The minimum number or percentage of available pods this budget requires. The name for the newly created object. The name for the newly created object. If not specified, the name of the input resource will be used. The name of the API generator to use. There are 2 generators: 'service/v1' and 'service/v2'. The only difference between them is that service port in v1 is named 'default', while it is left unnamed in v2. Default is 'service/v2'. The network protocol for the service to be created. Default is 'TCP'. The port that the service should serve on. Copied from the resource being exposed, if unspecified The resource requirement limits for this container.  For example, 'cpu=200m,memory=512Mi'.  Note that server side components may assign limits depending on the server configuration, such as limit ranges. The resource requirement requests for this container.  For example, 'cpu=100m,memory=256Mi'.  Note that server side components may assign requests depending on the server configuration, such as limit ranges. The type of secret to create Undo a previous rollout Update resource requests/limits on objects with pod templates Update the annotations on a resource Update the labels on a resource Update the taints on one or more nodes Username for Docker registry authentication View rollout history Where to output the files.  If empty or '-' uses stdout, otherwise creates a directory hierarchy in that directory dummy restart flag) kubectl controls the Kubernetes cluster manager Project-Id-Version: kubernetes
Report-Msgid-Bugs-To: EMAIL
PO-Revision-Date: 2017-08-28 15:20+0200
Last-Translator: Luca Berton <mr.evolution85@gmail.com>
Language-Team: Luca Berton <mr.evolution85@gmail.com>
Language: it_IT
MIME-Version: 1.0
Content-Type: text/plain; charset=UTF-8
Content-Transfer-Encoding: 8bit
Plural-Forms: nplurals=2; plural=(n != 1);
X-Generator: Poedit 1.8.7.1
X-Poedit-SourceCharset: UTF-8
 
		  # Mostra metriche per tutti i nodi
		 kubectl top node

		 # Mostra metriche per un determinato nodo
		 kubectl top node NODE_NAME 
		# Ottieni la documentazione della risorsa e i relativi campi
		kubectl explain pods

		# Ottieni la documentazione di un campo specifico di una risorsa
		kubectl explain pods.spec.containers 
		# Stampa i flag ereditati da tutti i comandi
		kubectl options 
		# Stampa le versioni client e server per il current context
		kubectl version 
		# Stampa le versioni API supportate
		kubectl api-versions 
		# Mostra metriche di tutti i pod nello spazio dei nomi predefinito
		kubectl top pod

		# Mostra metriche di tutti i pod nello spazio dei nomi specificato
		kubectl top pod --namespace=NAMESPACE

		# Mostra metriche per un pod e i suoi relativi container
		kubectl top pod POD_NAME --containers

		# Mostra metriche per i pod definiti da label name = myLabel
		kubectl top pod -l name=myLabel 
		Convertire i file di configurazione tra diverse versioni API. Sono
		accettati i formati YAML e JSON.

		Il comando prende il nome di file, la directory o l'URL come input e lo converte nel formato
		di versione specificata dal flag -output-version. Se la versione di destinazione non è specificata o
		non supportata, viene convertita nella versione più recente.

		L'output predefinito verrà stampato su stdout nel formato YAML. Si può usare l'opzione -o
		per cambiare la destinazione di output. 
		Creare un namespace con il nome specificato. 
		Crea un ruolo con una singola regola. 
		Creare un service account con il nome specificato. 
		Contrassegna il nodo come programmabile. 
		Contrassegnare il nodo come non programmabile. 
		Imposta le annotazioni dell'ultima-configurazione-applicata impostandola in modo che corrisponda al contenuto di un file.
		Ciò determina l'aggiornamento dell'ultima-configurazione-applicata come se 'kubectl apply -f <file>' fosse stato eseguito,
		senza aggiornare altre parti dell'oggetto. 
	  # Crea un nuovo namespace denominato my-namespace
	  kubectl create namespace my-namespace 
	  # Crea un nuovo service account denominato my-service-account
	  kubectl create serviceaccount my-service-account 
	Crea un servizio ExternalName con il nome specificato.

	Il servizio ExternalName fa riferimento a un indirizzo DNS esterno 
	solo pod, che permetteranno agli autori delle applicazioni di utilizzare i servizi di riferimento
	che esistono fuori dalla piattaforma, su altri cluster, o localmente.. 
	Help fornisce assistenza per qualsiasi comando nell'applicazione.
	Basta digitare kubectl help [path to command] per i dettagli completi. 
    # Creare un nuovo servizio LoadBalancer denominato my-lbs
    kubectl create service loadbalancer my-lbs --tcp=5678:8080 
    # Dump dello stato corrente del cluster verso stdout
    kubectl cluster-info dump

    # Dump dello stato corrente del cluster verso /path/to/cluster-state
    kubectl cluster-info dump --output-directory=/path/to/cluster-state

    # Dump di tutti i namespaces verso stdout
    kubectl cluster-info dump --all-namespaces

    # Dump di un set di namespace verso /path/to/cluster-state
    kubectl cluster-info dump --namespaces default,kube-system --output-directory=/path/to/cluster-state 
    Crea un servizio LoadBalancer con il nome specificato. Un insieme delimitato-da-virgole di quota scopes che devono corrispondere a ciascun oggetto gestito dalla quota. Un insieme delimitato-da-virgola di coppie risorsa = quantità che definiscono un hard limit. Un label selector da utilizzare per questo budget. Sono supportati solo i selettori equality-based selector. Un selettore di label da utilizzare per questo servizio. Sono supportati solo equality-based selector.  Se vuota (default) dedurre il selettore dal replication controller o replica set.) Indirizzo IP esterno aggiuntivo (non gestito da Kubernetes) da accettare per il servizio. Se questo IP viene indirizzato a un nodo, è possibile accedere da questo IP in aggiunta al service IP generato. Un override JSON inline per l'oggetto generato. Se questo non è vuoto, viene utilizzato per ignorare l'oggetto generato. Richiede che l'oggetto fornisca un campo valido apiVersion. Approva una richiesta di firma del certificato Assegnare il proprio ClusterIP o impostare su 'None' per un servizio 'headless' (nessun bilanciamento del carico). Collega a un container in esecuzione ClusterIP da assegnare al servizio. Lasciare vuoto per allocare automaticamente o impostare su 'None' per creare un servizio headless. ClusterRole a cui questo ClusterRoleBinding fa riferimento ClusterRole a cui questo RoleBinding fa riferimento Convertire i file di configurazione tra diverse versioni APIs Copiare file e directory da e verso i container. Crea un secret TLS Crea un namespace con il nome specificato Crea un secret da utilizzare con un registro Docker Crea un secret utilizzando un subcommand specificato Creare un account di servizio con il nome specificato Elimina il cluster specificato dal kubeconfig Elimina il context specificato dal kubeconfig Nega una richiesta di firma del certificato Descrive uno o più context Mostra i cluster definiti nel kubeconfig Visualizza le impostazioni merged di kubeconfig o un file kubeconfig specificato Visualizza una o più risorse Drain node in preparazione alla manutenzione Modificare una risorsa sul server Email per il registro Docker Esegui un comando in un contenitore Inoltra una o più porte locali a un pod Aiuto per qualsiasi comando Se non è vuoto, impostare l'affinità di sessione per il servizio; Valori validi: 'None', 'ClientIP' Se non è vuoto, l'aggiornamento delle annotazioni avrà successo solo se questa è la resource-version corrente per l'oggetto. Valido solo quando si specifica una singola risorsa. Se non vuoto, l'aggiornamento delle label avrà successo solo se questa è la resource-version corrente per l'oggetto. Valido solo quando si specifica una singola risorsa. Contrassegnare il nodo come programmabile Contrassegnare il nodo come non programmabile Imposta la risorsa indicata in pausa Modificare le risorse del certificato. Modifica i file kubeconfig Nome o numero di porta nel container verso il quale il servizio deve dirigere il traffico. Opzionale. Restituisce solo i log dopo una data specificata (RFC3339). Predefinito tutti i log. È possibile utilizzare solo uno tra data-inizio/a-partire-da. Codice di completamento shell di output per la shell specificata (bash o zsh) Password per l'autenticazione al registro di Docker Percorso certificato di chiave pubblica codificato PEM. Percorso alla chiave privata associata a un certificato specificato. Prerequisito per la versione delle risorse. Richiede che la versione corrente delle risorse corrisponda a questo valore per scalare. Stampa per client e server le informazioni sulla versione Stampa l'elenco flag ereditati da tutti i comandi Stampa i log per container in un pod Riprendere una risorsa in pausa Ruolo di riferimento per RoleBinding Esegui una particolare immagine nel cluster Eseguire un proxy al server Kubernetes API Posizione del server per il Registro Docker Imposta caratteristiche specifiche sugli oggetti Impostare il selettore di una risorsa Mostra i dettagli di una specifica risorsa o un gruppo di risorse Mostra lo stato del rollout Sinonimo di --target-port L'immagine per il container da eseguire. La politica di pull dell'immagine per il container. Se lasciato vuoto, questo valore non verrà specificato dal client e predefinito dal server Il numero minimo o la percentuale di pod disponibili che questo budget richiede. Il nome dell'oggetto appena creato. Il nome dell'oggetto appena creato. Se non specificato, verrà utilizzato il nome della risorsa di input. Il nome del generatore API da utilizzare. Ci sono 2 generatori: 'service/v1' e 'service/v2'. L'unica differenza tra loro è che la porta di servizio in v1 è denominata "predefinita", mentre viene lasciata unnamed in v2. Il valore predefinito è 'service/v2'. Il protocollo di rete per il servizio da creare. Il valore predefinito è 'TCP'. La porta che il servizio deve servire. Copiato dalla risorsa esposta, se non specificata I limiti delle richieste di risorse per questo contenitore.  Ad esempio, 'cpu=200m,memory=512Mi'. Si noti che i componenti lato server possono assegnare i limiti a seconda della configurazione del server, ad esempio intervalli di limiti. La risorsa necessita di richieste di requisiti per questo pod. Ad esempio, 'cpu = 100m, memoria = 256Mi'. Si noti che i componenti lato server possono assegnare i requisiti a seconda della configurazione del server, ad esempio intervalli di limiti. Tipo di segreto da creare Annulla un precedente rollout Aggiorna richieste di risorse/limiti sugli oggetti con pod template Aggiorna annotazioni di risorsa Aggiorna label di una risorsa Aggiorna i taints su uno o più nodi Nome utente per l'autenticazione nel registro Docker Visualizza la storia del rollout Dove eseguire l'output dei file. Se vuota o '-' utilizza lo stdout, altrimenti crea una gerarchia di directory in quella directory flag di riavvio finto) Kubectl controlla il gestore cluster di Kubernetes 