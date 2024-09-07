# home-k8s

# Secret
```
cd secret
echo -n 'XXX' | base64
(edit secret-origin.yaml)
kubeseal --format=yaml --cert=cert.pem < secret-origin.yaml > sealed-secret.yaml
```

