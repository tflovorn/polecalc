import json
import matplotlib.pyplot as plt
from matplotlib.ticker import FormatStrFormatter
from matplotlib.font_manager import FontProperties
from numpy import arange

def parse_file(file_path):
    '''Return the plot representation of the JSON file specified.'''
    # check for IOError?
    return import_json(open(file_path, 'r').read())

def import_json(json_string):
    '''Return the plot representation of the given JSON string.'''
    graph_data = json.loads(json_string)
    if isinstance(graph_data, list):
        graph_data = [_default_data(graph) for graph in graph_data]
    else:
        graph_data = _default_data(graph_data)
    return graph_data

def _default_data(graph_data):
    if "caption_font_size" not in graph_data:
        graph_data["caption_font_size"] = "large"

def make_graph(graph_data):
    '''Take a dictionary representing a graph or a list of such dictionaries.
    Build the graph(s), save them to file(s) (if requested), and return the
    matplotlib figures.

    '''
    if isinstance(graph_data, list):
        return [make_graph(some_graph) for some_graph in graph_data]
    fig = plt.figure()
    axes = fig.add_subplot(1, 1, 1)
    bounds = [None, None]
    for series in graph_data["series"]:
        fig, axes, bounds = _graph_series(graph_data, series, fig, axes, 
                                          bounds)
    fontprop = FontProperties(size=graph_data["caption_font_size"])
    

def _graph_series(graph_data, series, fig, axes, bounds):
    # -- todo : set ticks ---
    axes.plot(_xData(graph_data), _yData(graph_data), graph_data["style"], 
              label=graph_data["label"])
    return fig, axes, bounds
    

def _xData(series):
    return [point[0] for point in data]

def _yData(series):
    return [point[1] for point in data]
