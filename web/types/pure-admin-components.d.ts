declare module "jsbarcode" {
  const JsBarcode: any;
  export default JsBarcode;
}

declare module "cropperjs" {
  namespace Cropper {
    interface Options {
      [key: string]: any;
    }
  }

  class Cropper {
    constructor(element: any, options?: Cropper.Options);
    [key: string]: any;
  }

  export default Cropper;
}

declare module "qrcode" {
  export interface QRCodeRenderersOptions {
    [key: string]: any;
  }

  const QRCode: any;
  export default QRCode;
}

declare module "typeit" {
  const TypeIt: any;
  export type Options = Record<string, any>;
  export default TypeIt;
}

declare module "typeit/dist/types" {
  export type El = Element;
}

declare module "@amap/amap-jsapi-loader" {
  const AMapLoader: any;
  export default AMapLoader;
}

declare module "@logicflow/core" {
  export class LogicFlow {
    [key: string]: any;
  }

  const LogicFlowDefault: typeof LogicFlow;
  export default LogicFlowDefault;
}

declare module "vue-json-pretty" {
  const VueJsonPretty: any;
  export default VueJsonPretty;
}
