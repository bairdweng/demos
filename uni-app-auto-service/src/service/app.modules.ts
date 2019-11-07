
import { Module } from '@nestjs/common';
import { PackageModule } from './modules/package.module';


@Module({
    modules: [
        PackageModule
    ]
})

export class ApplicationModule { }