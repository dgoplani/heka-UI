export interface Notification {
  source: string,  // HOTFIX, NOTIFICATION
  type: string,
  tag: string,   
  title: string,
  body: string,
  action: string,
  extraData: any
}