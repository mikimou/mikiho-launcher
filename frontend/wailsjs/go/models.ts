export namespace main {
	
	export class Options {
	    nickname: string;
	    ram: number;
	
	    static createFrom(source: any = {}) {
	        return new Options(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.nickname = source["nickname"];
	        this.ram = source["ram"];
	    }
	}

}

