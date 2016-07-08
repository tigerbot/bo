(function () {
	'use strict';
	var gulp    = require('gulp')
	  ,	pug     = require('gulp-pug')
	  , browser = require('gulp-browser').browserify
	  ;

	gulp.task('html', function build_html() {
		return gulp.src('pug/*.pug')
			.pipe(pug({}))
			.pipe(gulp.dest('../public'))
		;
	});

	gulp.task('javascript', function combine_js() {
		return gulp.src('js/index.js')
			.pipe(browser())
			.pipe(gulp.dest('../public'))
		;
	});

	gulp.task('default', function () {
		gulp.start('html', 'javascript');
	});
	gulp.task('watch', function () {
		gulp.watch('pug/*.pug', ['html']);
		gulp.watch('js/*.js', ['javascript']);
	});
})();