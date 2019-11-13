
import { Module, MiddlewareConsumer, RequestMethod } from '@nestjs/common';
import { PackageController } from '../controllers/package.controller';
import { AppInfoController } from '../controllers/appinfo.controller';
@Module({
  controllers: [
    PackageController,
    AppInfoController
  ]
})
export class PackageModule {
  configure(consumer: MiddlewareConsumer) {
  }

  
}

