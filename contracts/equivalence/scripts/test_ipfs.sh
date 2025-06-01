# Teste manual de cada parte
echo "=== Teste 1: Download IPFS ==="
curl -L "http://127.0.0.1:8080/ipfs/QmT78zSuBmuS4z925WZfrqQ1qHaJ56DQaTfyMUF7F8ff5o"

echo -e "\n=== Teste 2: Validar JSON ==="
echo '{"title":"Teste","description":"ok"}' | jq .

echo -e "\n=== Teste 3: Contract Address ==="
echo "academic1jue5rlc9dkurt3etr57duutqu7prchqrk2mes2227m52kkrual3q82fkwv"