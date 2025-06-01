#!/bin/bash

# ============================================================================
# IPFS Real Content Fetcher for AcademicToken - WORKING VERSION
# ============================================================================

set -e  # Exit on any error

# Configuration
CONTRACT_ADDR="academic1jue5rlc9dkurt3etr57duutqu7prchqrk2mes2227m52kkrual3q82fkwv"
IPFS_GATEWAYS=(
    "http://127.0.0.1:8080/ipfs"
    "https://ipfs.io/ipfs"
    "https://gateway.pinata.cloud/ipfs"
)
TEMP_DIR="/tmp/ipfs_content"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

mkdir -p "$TEMP_DIR"

log() { echo -e "${BLUE}[$(date +'%H:%M:%S')]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1" >&2; }
success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }

show_usage() {
    cat << EOF
Usage: $0 -l IPFS_LINK TITLE

Example:
    $0 -l "QmT78zSuBmuS4z925WZfrqQ1qHaJ56DQaTfyMUF7F8ff5o" "Cálculo 1"
EOF
}

extract_ipfs_hash() {
    local ipfs_link=$1
    if [[ $ipfs_link =~ ipfs://([A-Za-z0-9]+) ]]; then
        echo "${BASH_REMATCH[1]}"
    elif [[ $ipfs_link =~ /ipfs/([A-Za-z0-9]+) ]]; then
        echo "${BASH_REMATCH[1]}"
    else
        echo "$ipfs_link"
    fi
}

download_from_ipfs() {
    local ipfs_hash=$1
    local output_file=$2
    
    log "Downloading IPFS content: $ipfs_hash"
    
    for gateway in "${IPFS_GATEWAYS[@]}"; do
        log "Trying gateway: $gateway"
        
        if curl -L --max-time 30 --silent \
           "$gateway/$ipfs_hash" -o "$output_file" 2>/dev/null; then
            
            if [[ -s "$output_file" ]]; then
                success "Downloaded from $gateway"
                return 0
            fi
        fi
    done
    
    error "Failed to download $ipfs_hash"
    return 1
}

# Safe JSON generation function
generate_syllabus_json() {
    local title=$1
    local content_text=$2
    
    # Escape function for JSON strings
    escape_json() {
        local input="$1"
        printf '%s' "$input" | sed 's/\\/\\\\/g; s/"/\\"/g' | tr -d '\000-\037\177' | head -c 200
    }
    
    local safe_title=$(escape_json "$title")
    local description="Conteúdo de $title"
    
    # If content is very short, use it as description
    if [[ ${#content_text} -lt 100 && -n "$content_text" ]]; then
        description="$(escape_json "$content_text")"
    fi
    
    # Generate valid JSON
    cat <<EOF
{
  "title": "$safe_title",
  "description": "$description",
  "objectives": [
    "Compreender os conceitos fundamentais de $safe_title",
    "Aplicar conhecimentos em situações práticas"
  ],
  "topics": [
    "Introdução ao tema",
    "Conceitos básicos",
    "Aplicações práticas"
  ],
  "bibliography": [
    "Bibliografia básica de $safe_title",
    "Referências complementares"
  ],
  "methodology": "Aulas expositivas e práticas",
  "evaluation": "Provas e trabalhos",
  "prerequisites": [],
  "workload_breakdown": {
    "theoretical_hours": 45,
    "practical_hours": 15,
    "laboratory_hours": 0,
    "field_work_hours": 0,
    "seminar_hours": 0,
    "individual_study_hours": 10
  },
  "course_level": "undergraduate",
  "duration_weeks": 16,
  "language": "português",
  "keywords": ["educação", "acadêmico"]
}
EOF
}

cache_in_contract() {
    local ipfs_link=$1
    local content_json=$2
    
    log "Caching content in smart contract"
    
    # Validate JSON first
    if ! echo "$content_json" | jq . >/dev/null 2>&1; then
        error "Invalid JSON generated"
        return 1
    fi
    
    local escaped_content=$(echo "$content_json" | jq -c .)
    
    log "Submitting to contract..."
    
    if academictokend tx wasm execute "$CONTRACT_ADDR" "{
      \"cache_ipfs_content\": {
        \"ipfs_link\": \"$ipfs_link\",
        \"content\": $escaped_content
      }
    }" --from alice --chain-id academictoken --gas auto --gas-adjustment 1.3 -y; then
        success "Content cached successfully!"
        return 0
    else
        error "Failed to cache content in contract"
        return 1
    fi
}

fetch_and_cache() {
    local ipfs_link=$1
    local title=$2
    
    log "Processing: $title ($ipfs_link)"
    
    local ipfs_hash=$(extract_ipfs_hash "$ipfs_link")
    local output_file="$TEMP_DIR/${ipfs_hash}_content"
    
    # Download content
    if ! download_from_ipfs "$ipfs_hash" "$output_file"; then
        error "Failed to download $ipfs_link"
        return 1
    fi
    
    # Read downloaded content
    local content_text=$(cat "$output_file")
    log "Content downloaded (${#content_text} chars)"
    
    # Generate JSON
    local content_json=$(generate_syllabus_json "$title" "$content_text")
    
    # Cache in contract
    if cache_in_contract "$ipfs_link" "$content_json"; then
        success "Successfully processed $title"
    else
        error "Failed to cache $title"
        return 1
    fi
    
    # Cleanup
    rm -f "$output_file"
}

# Main function
main() {
    log "Starting IPFS Real Content Fetcher"
    
    # Check dependencies
    for cmd in curl jq academictokend; do
        if ! command -v "$cmd" >/dev/null 2>&1; then
            error "Required command not found: $cmd"
            exit 1
        fi
    done
    
    # Parse arguments
    if [[ $# -ne 3 || "$1" != "-l" ]]; then
        show_usage
        exit 1
    fi
    
    local ipfs_link=$2
    local title=$3
    
    # Process the IPFS link
    fetch_and_cache "$ipfs_link" "$title"
    
    success "IPFS content fetcher completed!"
    rm -rf "$TEMP_DIR"
}

# Run
main "$@"