BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS users (
	id TEXT PRIMARY KEY,
	quota_max_instance INTEGER NOT NULL CHECK (quota_max_instance >= 0),
	quota_max_cpu INTEGER NOT NULL CHECK (quota_max_cpu >= 0),
	quota_max_memory INTEGER NOT NULL CHECK (quota_max_memory >= 0),
	status TEXT NOT NULL CHECK (status IN ('pending', 'initializing', 'active', 'failed')),
	error_phase TEXT CHECK (
		error_phase IS NULL OR
		error_phase IN ('failed in pending', 'failed in initializing')
	)
);

CREATE TABLE IF NOT EXISTS images (
	id TEXT PRIMARY KEY,
	alias TEXT NOT NULL UNIQUE,
	fingerprint TEXT NOT NULL UNIQUE,
	server_url TEXT NOT NULL,
	protocol TEXT NOT NULL,
	is_public INTEGER NOT NULL CHECK (is_public IN (0, 1))
);

CREATE TABLE IF NOT EXISTS vpcs (
	id TEXT PRIMARY KEY,
	owner_id TEXT NOT NULL UNIQUE,
	name TEXT NOT NULL,
	cidr TEXT NOT NULL UNIQUE,
	FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS subnets (
	id TEXT PRIMARY KEY,
	vpc_id TEXT NOT NULL,
	name TEXT NOT NULL,
	cidr TEXT NOT NULL,
	FOREIGN KEY (vpc_id) REFERENCES vpcs(id) ON DELETE CASCADE,
	UNIQUE (vpc_id, name),
	UNIQUE (vpc_id, cidr),
	UNIQUE (cidr)
);

CREATE TABLE IF NOT EXISTS leases (
	subnet_id TEXT NOT NULL,
	target_id TEXT PRIMARY KEY,
	ip_address TEXT NOT NULL,
	FOREIGN KEY (subnet_id) REFERENCES subnets(id) ON DELETE CASCADE,
	UNIQUE (subnet_id, ip_address)
);

CREATE TABLE IF NOT EXISTS volumes (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	size_gb INTEGER NOT NULL CHECK (size_gb > 0),
	pool TEXT NOT NULL,
	owner TEXT NOT NULL,
	FOREIGN KEY (owner) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS instances (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	owner_id TEXT NOT NULL,
	status TEXT NOT NULL CHECK (
		status IN (
			'pending', 'creating', 'starting', 'running',
			'stopping', 'stopped', 'deleting', 'error'
		)
	),
	error_phase TEXT CHECK (
		error_phase IS NULL OR
		error_phase IN (
			'error in pending',
			'error in creating',
			'error in starting',
			'error in stopping',
			'error in deleting'
		)
	),
	cpu INTEGER NOT NULL CHECK (cpu > 0),
	memory_mb INTEGER NOT NULL CHECK (memory_mb > 0),
	image_id TEXT NOT NULL,
	subnet_id TEXT NOT NULL,
	private_ip TEXT NOT NULL,
	root_volume_id TEXT NOT NULL UNIQUE,
	FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE RESTRICT,
	FOREIGN KEY (subnet_id) REFERENCES subnets(id) ON DELETE RESTRICT,
	FOREIGN KEY (root_volume_id) REFERENCES volumes(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS ingress_routes (
	id TEXT PRIMARY KEY,
	subdomain TEXT NOT NULL,
	port_name TEXT NOT NULL,
	target_ip TEXT NOT NULL,
	target_port INTEGER NOT NULL CHECK (target_port > 0 AND target_port <= 65535),
	instance_id TEXT NOT NULL,
	owner_id TEXT NOT NULL,
	FOREIGN KEY (instance_id) REFERENCES instances(id) ON DELETE CASCADE,
	FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE,
	UNIQUE (subdomain, port_name)
);

CREATE INDEX IF NOT EXISTS idx_instances_owner_id ON instances(owner_id);
CREATE INDEX IF NOT EXISTS idx_instances_subnet_id ON instances(subnet_id);
CREATE INDEX IF NOT EXISTS idx_ingress_routes_instance_id ON ingress_routes(instance_id);
CREATE INDEX IF NOT EXISTS idx_ingress_routes_owner_id ON ingress_routes(owner_id);
CREATE INDEX IF NOT EXISTS idx_subnets_vpc_id ON subnets(vpc_id);
CREATE INDEX IF NOT EXISTS idx_leases_subnet_id ON leases(subnet_id);

COMMIT;
