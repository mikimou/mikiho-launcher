export namespace main {
	
	export class ModpackManifest {
	    version: string;
	    url: string;
	    command?: string[];
	
	    static createFrom(source: any = {}) {
	        return new ModpackManifest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.url = source["url"];
	        this.command = source["command"];
	    }
	}
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

