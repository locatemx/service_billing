# service_billing

#ACCOUNT#

id uuid
parent uuid
billing character varying
wallet character varying
type character varying
name character varying
live boolean
active boolean
deleted boolean
created_at timestamp with time zone
updated_at timestamp with time zone
account character varying
credit boolean

#ACCOUNT_DEVICE#

id uuid
account uuid
device uuid
service uuid
aditional uuid
alias character varying
props jsonb
created_at timestamp with time zone

#BILLING_SERVICE#

id uuid
account uuid
base uuid
alias character varying
price numeric
deleted boolean
created_at timestamp with time zone

#BILLING_SERVICE_BASE#

id uuid
name character varying
description text
invoice_name character varying
default_price numeric
max_price numeric
min_price numeric
interval character varying
interval_count smallint
account_limit smallint
is_aditional boolean
retention smallint
next_level uuid
fee numeric

#BILLING_SERVICE#

id uuid
account uuid
base uuid
alias character varying
price numeric
deleted boolean
created_at timestamp with time zone

#BILLING_SERVICE_BASE#

id uuid
name character varying
description text
invoice_name character varying
default_price numeric
min_price numeric
max_price numeric

due_date
external_reference
