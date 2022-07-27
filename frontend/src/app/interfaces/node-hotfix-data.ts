export interface NodeHotfixData {
  unique_id: string
  ip: string
  hostname: string
  role: string
  ha_enable: boolean
  status: string
  master_candidate: boolean
  hotfixes: Hotfix[]
}

interface Hotfix {
  name: string
  timestamp: string
  status: string
}
  
