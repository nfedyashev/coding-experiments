# README

This is just an experiment in isomorophic validation logic sharing between JS & ruby models

    irb(main):001:0> Article.new(title: 'boom').valid?
    => true
    irb(main):002:0> Article.new(title: 'bo balalaika om').valid?
    => false

# TODO

 * Use precompliled version instead instead of instantiating V8 context on each validator call, see https://github.com/cowboyd/therubyracer/blob/master/spec/c/script_spec.rb
 * Properly export vanilla validator js function so that it could be called in real validation library(like jquery.validate)
 * Pass error message along with exported function to properly match them(cxt object may contain all exported functions/objects)



# NOTES

Example usage of Custom jquery.validator logic

    jQuery.validator.addMethod("greaterThanZero", function(value, element) {
        return this.optional(element) || (parseFloat(value) > 0);
    }, "* Amount must be greater than zero");
