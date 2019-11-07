
import { Module, MiddlewaresConsumer, RequestMethod } from '@nestjs/common';
import { PackageController } from '../controllers/package.controller';



@Module({
  controllers: [
    PackageController
  ]
})
export class PackageModule {
  configure(consumer: MiddlewaresConsumer) {
  }
}

