// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'label_state.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

T _$identity<T>(T value) => value;

final _privateConstructorUsedError = UnsupportedError(
  'It seems like you constructed your class using `MyClass._()`. This constructor is only meant to be used by freezed and you are not supposed to need it nor use it.\nPlease check the documentation here for more information: https://github.com/rrousselGit/freezed#adding-getters-and-methods-to-our-models',
);

/// @nodoc
mixin _$LabelState {
  List<LabelItem> get labels => throw _privateConstructorUsedError;
  bool get isLoading => throw _privateConstructorUsedError;

  /// Create a copy of LabelState
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $LabelStateCopyWith<LabelState> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $LabelStateCopyWith<$Res> {
  factory $LabelStateCopyWith(
    LabelState value,
    $Res Function(LabelState) then,
  ) = _$LabelStateCopyWithImpl<$Res, LabelState>;
  @useResult
  $Res call({List<LabelItem> labels, bool isLoading});
}

/// @nodoc
class _$LabelStateCopyWithImpl<$Res, $Val extends LabelState>
    implements $LabelStateCopyWith<$Res> {
  _$LabelStateCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of LabelState
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({Object? labels = null, Object? isLoading = null}) {
    return _then(
      _value.copyWith(
            labels:
                null == labels
                    ? _value.labels
                    : labels // ignore: cast_nullable_to_non_nullable
                        as List<LabelItem>,
            isLoading:
                null == isLoading
                    ? _value.isLoading
                    : isLoading // ignore: cast_nullable_to_non_nullable
                        as bool,
          )
          as $Val,
    );
  }
}

/// @nodoc
abstract class _$$LabelStateImplCopyWith<$Res>
    implements $LabelStateCopyWith<$Res> {
  factory _$$LabelStateImplCopyWith(
    _$LabelStateImpl value,
    $Res Function(_$LabelStateImpl) then,
  ) = __$$LabelStateImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({List<LabelItem> labels, bool isLoading});
}

/// @nodoc
class __$$LabelStateImplCopyWithImpl<$Res>
    extends _$LabelStateCopyWithImpl<$Res, _$LabelStateImpl>
    implements _$$LabelStateImplCopyWith<$Res> {
  __$$LabelStateImplCopyWithImpl(
    _$LabelStateImpl _value,
    $Res Function(_$LabelStateImpl) _then,
  ) : super(_value, _then);

  /// Create a copy of LabelState
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({Object? labels = null, Object? isLoading = null}) {
    return _then(
      _$LabelStateImpl(
        labels:
            null == labels
                ? _value._labels
                : labels // ignore: cast_nullable_to_non_nullable
                    as List<LabelItem>,
        isLoading:
            null == isLoading
                ? _value.isLoading
                : isLoading // ignore: cast_nullable_to_non_nullable
                    as bool,
      ),
    );
  }
}

/// @nodoc

class _$LabelStateImpl implements _LabelState {
  const _$LabelStateImpl({
    required final List<LabelItem> labels,
    this.isLoading = false,
  }) : _labels = labels;

  final List<LabelItem> _labels;
  @override
  List<LabelItem> get labels {
    if (_labels is EqualUnmodifiableListView) return _labels;
    // ignore: implicit_dynamic_type
    return EqualUnmodifiableListView(_labels);
  }

  @override
  @JsonKey()
  final bool isLoading;

  @override
  String toString() {
    return 'LabelState(labels: $labels, isLoading: $isLoading)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$LabelStateImpl &&
            const DeepCollectionEquality().equals(other._labels, _labels) &&
            (identical(other.isLoading, isLoading) ||
                other.isLoading == isLoading));
  }

  @override
  int get hashCode => Object.hash(
    runtimeType,
    const DeepCollectionEquality().hash(_labels),
    isLoading,
  );

  /// Create a copy of LabelState
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$LabelStateImplCopyWith<_$LabelStateImpl> get copyWith =>
      __$$LabelStateImplCopyWithImpl<_$LabelStateImpl>(this, _$identity);
}

abstract class _LabelState implements LabelState {
  const factory _LabelState({
    required final List<LabelItem> labels,
    final bool isLoading,
  }) = _$LabelStateImpl;

  @override
  List<LabelItem> get labels;
  @override
  bool get isLoading;

  /// Create a copy of LabelState
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$LabelStateImplCopyWith<_$LabelStateImpl> get copyWith =>
      throw _privateConstructorUsedError;
}
