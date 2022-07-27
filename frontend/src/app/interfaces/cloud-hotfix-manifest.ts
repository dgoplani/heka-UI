export interface CloudHotfixManifest {
  type: string
  metadata: CloudHotfixMetadata
  data: Hotfix[]
}

export interface CloudHotfixMetadata {
  version: number
  generated: string
}

export interface Hotfix {
  name: string
  sha256: string
  revert: HotfixRevert
  ticketId: string
  released: string
  compatibleReleases: string[]
  type: string
  compatibleNode: string
  summary: string
  impactedArea: string[]
  fixes: HotfixFixes
  severity: string
  references: HotfixReference[]
  requiredActions: HotfixRequiredActions
  incompatible: string[]
}

export interface ProcessedHotfix {
  name: string
  sha256: string
  revert: HotfixRevert
  ticketId: string
  released: string
  compatibleReleases: string[]
  type: string
  compatibleNode: string
  summary: string
  impactedArea: string[]
  fixes: HotfixFixes
  severity: string
  references: HotfixReference[]
  requiredActions: HotfixRequiredActions
  incompatible: string[]
  applyStatus: string
  applyTimestamp: string
  show: boolean
  idx: number
}

export interface HotfixRevert {
  name: string
  sha256: string
}

export interface HotfixFixes {
  BUGFIX: HotfixFixesBugfix[]
  CVE: HotfixFixesCveFix[]
  SECURITY: any[]
}

export interface HotfixFixesBugfix {
  id: string
  summary: string
}

export interface HotfixFixesCveFix {
  id: string
  summary: string
  ref: string
}

export interface HotfixReference {
  type: string
  link: string
}

export interface HotfixRequiredActions {
  systemReboot: string
  productRestart: string
  serviceRestart: string[]
}
      
