export namespace main {
	
	export class FormData {
	    sourceDir: string;
	    includeChild: boolean;
	    targetDir: string;
	    dirType: string;
	    isMove: boolean;
	    noShotTimeType: number;
	    skipSameFile: boolean;
	    skipFileLessThan: number;
	    skipFileContains: string;
	
	    static createFrom(source: any = {}) {
	        return new FormData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sourceDir = source["sourceDir"];
	        this.includeChild = source["includeChild"];
	        this.targetDir = source["targetDir"];
	        this.dirType = source["dirType"];
	        this.isMove = source["isMove"];
	        this.noShotTimeType = source["noShotTimeType"];
	        this.skipSameFile = source["skipSameFile"];
	        this.skipFileLessThan = source["skipFileLessThan"];
	        this.skipFileContains = source["skipFileContains"];
	    }
	}

}

