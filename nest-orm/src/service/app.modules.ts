
import { Module } from '@nestjs/common';
import { PackageModule } from './modules/package.module';


@Module({
    imports:[
        PackageModule
    ]
})

export class ApplicationModule { }