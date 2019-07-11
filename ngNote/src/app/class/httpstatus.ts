export class HTTPStatus {
    Status: HttpStatusCode;
    Message: any;
    Error: any;
}

export enum HttpStatusCode {
    Ok = 'ok',
    NotFound = 'not_found',
    Error = 'error'
}