export namespace organizer {
	
	export class ProgressEvent {
	    total: number;
	    processed: number;
	    moved: number;
	    skipped: number;
	    errors: number;
	    currentFile: string;
	    percentDone: number;
	    categoryCounts: Record<string, number>;
	    running: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ProgressEvent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.processed = source["processed"];
	        this.moved = source["moved"];
	        this.skipped = source["skipped"];
	        this.errors = source["errors"];
	        this.currentFile = source["currentFile"];
	        this.percentDone = source["percentDone"];
	        this.categoryCounts = source["categoryCounts"];
	        this.running = source["running"];
	    }
	}
	export class Summary {
	    total: number;
	    moved: number;
	    skipped: number;
	    errors: number;
	    elapsedSeconds: number;
	    categoryCounts: Record<string, number>;
	    logPath: string;
	    dryRun: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Summary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.moved = source["moved"];
	        this.skipped = source["skipped"];
	        this.errors = source["errors"];
	        this.elapsedSeconds = source["elapsedSeconds"];
	        this.categoryCounts = source["categoryCounts"];
	        this.logPath = source["logPath"];
	        this.dryRun = source["dryRun"];
	    }
	}

}

