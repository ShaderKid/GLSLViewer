#version 410 core

out vec4 Color;

uniform float time;
uniform vec2 resolution;
uniform vec2 mouse;

void main() {
    //Color = vec4(sin(time),0,1,1);
    Color = vec4(mouse.x,mouse.y,1,1);
}
