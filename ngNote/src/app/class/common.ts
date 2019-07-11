import { Injector } from '@angular/core';

export let InjectorInstance: Injector;
export const ApiUrl: string = "http://localhost:8080/api";

export function InitInjector(injector: Injector) {
    if (InjectorInstance == null)
        InjectorInstance = injector
}
