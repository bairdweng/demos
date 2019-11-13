import { AppInfo } from "../entitys/appinfo";
import { getConnection } from "typeorm";

export class appInfoService {
  async getInfo(): Promise<AppInfo> {
    const appInfo = await getConnection()
      .createQueryBuilder()
      .select("app_info")
      .from(AppInfo, "app_info")
      .where("app_info.id = :id", { id: 1 })
      .getOne();
    return appInfo
  }
  
}