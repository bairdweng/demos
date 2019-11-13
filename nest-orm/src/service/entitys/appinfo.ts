
import { Entity, Column, PrimaryGeneratedColumn } from "typeorm";
// 可映射数据表。会自动忘数据库同步表结构，谨慎删除字段。
@Entity("app_info")
export class AppInfo {
  @PrimaryGeneratedColumn()
  id: number;

  @Column()
  name:string;

  @Column()
  bundle_id: string;

  @Column()
  key_name:String;
  
  @Column()
  profile_name:String;
  
  @Column()
  images:number;
  
  @Column()
  status:number;

  @Column()
  create_time:Date;
  
  @Column()
  update_time:Date;
}