import sys
import json
import matplotlib.pyplot as plt
from matplotlib.ticker import FormatStrFormatter
from matplotlib.font_manager import FontProperties
from numpy import arange

_GRAPH_DEFAULTS = {"xlabel":"$x$", "ylabel":"$y$", "num_ticks":5, 
    "axis_label_fontsize":"large", "tick_formatstr":"%.2f",
    "legend_fontsize":"large", "legend_loc":0, "legend_title":None, 
    "graph_filepath":None}

_SERIES_DEFAULTS = {"label":None, "style":"k."}

def parse_file(file_path):
    '''Return the plot representation of the JSON file specified.'''
    # -- todo : check for IOError --
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
    # graph-wide defaults
    for key, value in _GRAPH_DEFAULTS.items():
        if key not in graph_data:
            graph_data[key] = value
    # hack to give a fresh series list each time
    if "series" not in graph_data:
        graph_data["series"] = []
    # series-specific defaults
    for series in graph_data["series"]:
        for key, value in _SERIES_DEFAULTS.items():
            if key not in series:
                series[key] = value
    return graph_data

def make_graph(graph_data):
    '''Take a dictionary representing a graph or a list of such dictionaries.
    Build the graph(s), save them to file(s) (if requested), and return the
    matplotlib figures.

    '''
    if isinstance(graph_data, list):
        return [make_graph(some_graph) for some_graph in graph_data]
    try:
        dims = graph_data["dimensions"]
        fig = plt.figure(figsize=(dims[0], dims[1]))
    except:
        fig = plt.figure()
    axes = fig.add_subplot(1, 1, 1)
    bounds = [None, None]
    for series in graph_data["series"]:
        fig, axes, bounds = _graph_series(graph_data, series, fig, axes, 
                                          bounds)
    fontprop = FontProperties(size=graph_data["legend_fontsize"])
    axes.legend(loc=graph_data["legend_loc"], title=graph_data["legend_title"],
                prop=fontprop)
    _save_figure(graph_data, fig)
    return fig, axes
    

def _graph_series(graph_data, series, fig, axes, bounds):
    # -- todo : set ticks --
    axes.plot(_xData(series), _yData(series), series["style"], 
              label=series["label"])
    return fig, axes, bounds

def _xData(series):
    return [point[0] for point in series["data"]]

def _yData(series):
    return [point[1] for point in series["data"]]

def _save_figure(graph_data, fig):
    if graph_data["graph_filepath"] is None:
        return
    fig.savefig(graph_data["graph_filepath"] + ".png")
    fig.savefig(graph_data["graph_filepath"] + ".eps")

if __name__ == "__main__":
    if len(sys.argv) > 1:
        make_graph(parse_file(sys.argv[1]))
